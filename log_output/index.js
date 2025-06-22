const express = require('express');
const { v4: uuidv4 } = require('uuid');

const randomString = uuidv4();

const app = express();
const PORT = 3000;

app.get('/', (req, res) => {
  const currentTimestamp = new Date();

  res.json({
    timestamp: currentTimestamp.toISOString(),
    randomString: randomString
  });
});

function printToConsole() {
    const timestamp = new Date().toISOString();
    console.log(`${timestamp}: ${randomString}`);
}

setInterval(printToConsole, 5000);

app.listen(PORT, () => {
  console.log(`Server is running and listening on http://localhost:${PORT}`);
  printToConsole();
});