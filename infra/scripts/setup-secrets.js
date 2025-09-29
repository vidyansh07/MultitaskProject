#!/usr/bin/env node

/**
 * Setup AWS SSM Parameters for different environments
 * This script helps initialize the required secrets in AWS Parameter Store
 */

const AWS = require('aws-sdk');
const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

const question = (prompt) => {
  return new Promise((resolve) => {
    rl.question(prompt, resolve);
  });
};

const PARAMETERS = [
  {
    name: 'jwt-secret',
    description: 'JWT Secret for token signing',
    secure: true,
    defaultValue: () => require('crypto').randomBytes(32).toString('hex')
  },
  {
    name: 'gemini-api-key',
    description: 'Google Gemini API Key',
    secure: true,
    defaultValue: null
  },
  {
    name: 'openai-api-key',
    description: 'OpenAI API Key',
    secure: true,
    defaultValue: null
  }
];

async function setupParameters() {
  console.log('🚀 Setting up AWS SSM Parameters for Multitask Platform\n');

  const region = await question('Enter AWS region (default: us-east-1): ') || 'us-east-1';
  const stage = await question('Enter stage (dev/staging/prod): ');
  
  if (!['dev', 'staging', 'prod'].includes(stage)) {
    console.error('❌ Invalid stage. Must be dev, staging, or prod');
    process.exit(1);
  }

  const ssm = new AWS.SSM({ region });

  console.log(`\n📍 Setting up parameters for stage: ${stage} in region: ${region}\n`);

  for (const param of PARAMETERS) {
    const parameterName = `/multitask/${stage}/${param.name}`;
    
    try {
      // Check if parameter already exists
      try {
        await ssm.getParameter({ Name: parameterName }).promise();
        console.log(`✅ Parameter ${parameterName} already exists`);
        continue;
      } catch (err) {
        if (err.code !== 'ParameterNotFound') {
          throw err;
        }
      }

      let value = param.defaultValue ? param.defaultValue() : null;
      
      if (!value) {
        value = await question(`Enter value for ${param.description}: `);
      }

      if (!value) {
        console.log(`⚠️  Skipping ${parameterName} (no value provided)`);
        continue;
      }

      await ssm.putParameter({
        Name: parameterName,
        Value: value,
        Type: param.secure ? 'SecureString' : 'String',
        Description: param.description,
        Tags: [
          { Key: 'Project', Value: 'multitask-platform' },
          { Key: 'Environment', Value: stage },
          { Key: 'ManagedBy', Value: 'setup-script' }
        ]
      }).promise();

      console.log(`✅ Created parameter: ${parameterName}`);
    } catch (error) {
      console.error(`❌ Failed to create parameter ${parameterName}:`, error.message);
    }
  }

  console.log('\n🎉 Parameter setup completed!');
  console.log('\n📚 Next steps:');
  console.log('1. Review the parameters in AWS Systems Manager Parameter Store');
  console.log('2. Update any placeholder values with real API keys');
  console.log('3. Run: npm run deploy:' + stage);
  
  rl.close();
}

if (require.main === module) {
  setupParameters().catch(console.error);
}

module.exports = { setupParameters };