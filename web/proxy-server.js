const http = require('http');
const url = require('url');
const path = require('path');
const fs = require('fs');

const FRONTEND_PORT = process.env.FRONTEND_PORT || 30005;
const BACKEND_PORT = process.env.BACKEND_PORT || 30002;
const BACKEND_HOST = process.env.BACKEND_HOST || '127.0.0.1';

const mimeTypes = {
    '.html': 'text/html',
    '.js': 'text/javascript',
    '.css': 'text/css',
    '.json': 'application/json',
    '.png': 'image/png',
    '.jpg': 'image/jpeg',
    '.gif': 'image/gif',
    '.svg': 'image/svg+xml',
    '.ico': 'image/x-icon',
    '.woff': 'application/font-woff',
    '.woff2': 'application/font-woff2',
    '.ttf': 'application/font-ttf',
    '.eot': 'application/vnd.ms-fontobject',
    '.map': 'application/json'
};

function serveStatic(req, res) {
    let filePath = path.join(__dirname, 'dist', req.url === '/' ? 'index.html' : req.url);

    const ext = path.extname(filePath);
    const contentType = mimeTypes[ext] || 'application/octet-stream';

    fs.readFile(filePath, (err, content) => {
        if (err) {
            if (err.code === 'ENOENT') {
                fs.readFile(path.join(__dirname, 'dist', 'index.html'), (err, content) => {
                    res.writeHead(200, { 'Content-Type': 'text/html' });
                    res.end(content, 'utf-8');
                });
            } else {
                res.writeHead(500);
                res.end('Server Error');
            }
        } else {
            res.writeHead(200, { 'Content-Type': contentType });
            res.end(content, 'utf-8');
        }
    });
}

function proxyAPI(req, res) {
    const targetUrl = `http://${BACKEND_HOST}:${BACKEND_PORT}${req.url}`;

    const options = {
        hostname: BACKEND_HOST,
        port: BACKEND_PORT,
        path: req.url,
        method: req.method,
        headers: {
            ...req.headers,
            'X-Forwarded-For': req.connection.remoteAddress,
            'X-Forwarded-Proto': 'http'
        }
    };

    const proxy = http.request(options, (response) => {
        res.writeHead(response.statusCode, response.headers);
        response.pipe(res, { end: true });
    });

    req.pipe(proxy, { end: true });

    proxy.on('error', (err) => {
        console.error('Proxy error:', err.message);
        res.writeHead(502, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ code: 502, message: 'Backend unavailable' }));
    });
}

const server = http.createServer((req, res) => {
    const parsedUrl = url.parse(req.url, true);

    if (req.url.startsWith('/api')) {
        proxyAPI(req, res);
    } else {
        serveStatic(req, res);
    }
});

server.listen(FRONTEND_PORT, () => {
    console.log(`[PROXY] Server running at http://0.0.0.0:${FRONTEND_PORT}/`);
    console.log(`[PROXY] Proxying /api/* to http://${BACKEND_HOST}:${BACKEND_PORT}`);
});
