const express = require('express');
const fs = require('fs').promises;
const path = require('path');

const PORT = 3000;
const LOG_FILE_PATH = path.join(__dirname, 'log', 'app_log.txt');
const DATA_FILE_PATH = path.join(__dirname, 'data', 'pingpong_count.txt');

const app = express();

async function getLastLine(filePath) {
  try {
    const fileHandle = await fs.open(filePath, 'r');
    const stats = await fileHandle.stat();

    if (stats.size === 0) {
      await fileHandle.close();
      return '';
    }

    let lastLine = '';
    let buffer = Buffer.alloc(1);
    for (let i = stats.size - 1; i >= 0; i--) {
      await fileHandle.read(buffer, 0, 1, i);
      const char = buffer.toString('utf8');

      if (char === '\n' && lastLine.length > 0) {
        break;
      }
      lastLine = char + lastLine;
    }

    await fileHandle.close();
    return lastLine.trim();

  } catch (err) {
    if (err.code === 'ENOENT') {
      return 'Log file not yet created.';
    }
    throw err;
  }
}

app.get('/', async (req, res) => {
  try {
    const logFilePromise = getLastLine(LOG_FILE_PATH)

    const dataFilePromise = fs.readFile(DATA_FILE_PATH, 'utf8')
      .catch(err => {
        if (err.code === 'ENOENT') return '0';
        throw err;
      });

    const [logContent, dataContent] = await Promise.all([
      logFilePromise,
      dataFilePromise,
    ]);

    const responseText = `${logContent}\nPing / Pongs: ${dataContent}`;

    res.status(200).type('text/plain').send(responseText);

  } catch (error) {
    console.error('An unexpected error occurred:', error);
    res.status(500).type('text/plain').send('An internal server error occurred.');
  }
});

app.listen(PORT, async () => {
  console.log(`Listening for requests on http://localhost:${PORT}`);
});