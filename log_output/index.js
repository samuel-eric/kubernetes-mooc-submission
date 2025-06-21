const { v4: uuidv4 } = require('uuid');

const randomString = uuidv4();

function printToConsole(randomString) {
    const timestamp = new Date().toISOString();
    console.log(`${timestamp}: ${randomString}`);
}

printToConsole(randomString)

setInterval(() => {
  printToConsole(randomString)
}, 5000);
