const fs = require('fs');
const path = require('path');
const express = require('express');

const DATA_DIR = path.join(__dirname, 'data');
const DATA_FILE_PATH = path.join(DATA_DIR, 'pingpong_count.txt');
const app = express();
const PORT = 3001;

if (!fs.existsSync(DATA_DIR)) {
  fs.mkdirSync(DATA_DIR);
}

app.get('/pingpong', (req, res) => {
  let currentCounter = 0;

  fs.readFile(DATA_FILE_PATH, 'utf8', (readErr, data) => {
    if (readErr && readErr.code !== 'ENOENT') {
      console.error('Failed to read file:', readErr);
      return res.status(500).send('Internal server error on read');
    }

    if (!readErr) {
      currentCounter = parseInt(data, 10) || 0;
    }

    currentCounter++;
    const newCounterValue = currentCounter;

    fs.writeFile(DATA_FILE_PATH, `${newCounterValue}`, (writeErr) => {
      if (writeErr) {
        console.error('Failed to write to file:', writeErr);
        return res.status(500).send('Internal server error on write');
      }

      res.status(200).send(`pong ${newCounterValue}`);
    });
  });
});

app.listen(PORT, () => {
  console.log(`Ping-pong server is running on port ${PORT}`);
});