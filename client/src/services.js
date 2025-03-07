import path from 'path';
import {spawn} from 'child_process';
import os from 'os';
import {app} from 'electron';
import express from 'express';

const DEBUG = process.env.NODE_ENV === 'development';

let _ztfServerProcess;

export function startZtfServer() {
    if (process.env.SKIP_SERVER) {
        console.log(`>> Skip to start ZTF Server by env "SKIP_SERVER=${process.env.SKIP_SERVER}".`);
        return Promise.resolve();
    }
    if (_ztfServerProcess) {
        return Promise.resolve(_ztfServerProcess);
    }

    let {SERVER_EXE_PATH: serverExePath} = process.env;
    if (!serverExePath && !DEBUG) {
        const platform = os.platform(); // 'darwin', 'linux', 'win32'
        const exePath = `bin/${platform}/ztf${platform === 'win32' ? '.exe' : ''}`;
        serverExePath = path.resolve(process.resourcesPath, exePath);
    }
    if (serverExePath) {
        if (!path.isAbsolute(serverExePath)) {
            serverExePath = path.resolve(app.getAppPath(), serverExePath);
        }
        return new Promise((resolve, reject) => {
            const cwd = process.env.SERVER_CWD_PATH || path.dirname(serverExePath);
            console.log(`>> Starting ZTF Server from exe path with command "${serverExePath} -P 8085" in "${cwd}"...`);
            const cmd = spawn(serverExePath, ['-P', '8085'], {
                cwd,
                shell: true,
            });
            cmd.on('close', (code) => {
                console.log(`>> ZTF server closed with code ${code}`);
                _ztfServerProcess = null;
            });
            cmd.stdout.on('data', data => {
                const dataString = String(data);
                const lines = dataString.split('\n');
                for (let i = 0; i < lines.length; i++) {
                    const line = lines[i];
                    if (DEBUG) {
                        console.log('\t', line);
                    }
                    if (line.includes('Now listening on: http')) {
                        resolve(line.split('Now listening on:')[1].trim());
                        if (!DEBUG) {
                            break;
                        }
                    } else if (line.includes('启动HTTP服务于')) {
                        resolve(line.split(/启动HTTP服务于|，/)[1].trim());
                        if (!DEBUG) {
                            break;
                        }
                    } else if (line.startsWith('[ERRO]')) {
                        reject(new Error(`Start ztf server failed with error: ${line.substring('[ERRO]'.length)}`));
                        if (!DEBUG) {
                            break;
                        }
                    }
                }
            });
            cmd.on('error', spawnError => {
                console.error('>>> Start ztf server failed with error', spawnError);
                reject(spawnError)
            });
            _ztfServerProcess = cmd;
        });
    }

    return new Promise((resolve, reject) => {
        const cwd = process.env.SERVER_CWD_PATH || path.resolve(app.getAppPath(), '../');
        console.log(`>> Starting ZTF development server from source with command "go run cmd/server/main.go -P 8085" in "${cwd}"`);
        const cmd = spawn('go', ['run', 'main.go', '-P', '8085'], {
            cwd,
            shell: true,
        });
        cmd.on('close', (code) => {
            console.log(`>> ZTF server closed with code ${code}`);
            _ztfServerProcess = null;
        });
        cmd.stdout.on('data', data => {
            const dataString = String(data);
            const lines = dataString.split('\n');
            for (let i = 0; i < lines.length; i++) {
                const line = lines[i];
                if (DEBUG) {
                    console.log('\t', line);
                }
                if (line.includes('Now listening on: http')) {
                    resolve(line.split('Now listening on:')[1].trim());
                    if (!DEBUG) {
                        break;
                    }
                } else if (line.startsWith('[ERRO]')) {
                    reject(new Error(`Start ztf server failed with error: ${line.substring('[ERRO]'.length)}`));
                    if (!DEBUG) {
                        break;
                    }
                }
            }
        });
        cmd.on('error', spawnError => {
            console.error('>>> Start ztf server failed with error', spawnError);
            reject(spawnError)
        });
        _ztfServerProcess = cmd;
    });
}

let _uiServerApp;

export function getUIServerUrl() {
    if (_uiServerApp) {
        return Promise.resolve();
    }

    let {UI_SERVER_URL: uiServerUrl} = process.env;
    if (!uiServerUrl && !DEBUG) {
        uiServerUrl = path.resolve(process.resourcesPath, 'ui');
    }

    if (uiServerUrl) {
        if (/^https?:\/\//.test(uiServerUrl)) {
            return Promise.resolve(uiServerUrl);
        }
        return new Promise((resolve, reject) => {
            if (!path.isAbsolute(uiServerUrl)) {
                uiServerUrl = path.resolve(app.getAppPath(), uiServerUrl);
            }

            const port = process.env.UI_SERVER_PORT || 8000;
            console.log(`>> Starting UI serer at ${uiServerUrl} with port ${port}`);

            const uiServer = express();
            uiServer.use(express.static(uiServerUrl));
            const server = uiServer.listen(port, serverError => {
                if (serverError) {
                    console.error('>>> Start ui server failed with error', serverError);
                    _uiServerApp = null;
                    reject(serverError);
                } else {
                    console.log(`>> UI server started successfully on http://localhost:${port}.`);
                    resolve(`http://localhost:${port}`);
                }
            });
            server.on('close', () => {
                _uiServerApp = null;
            });
            _uiServerApp = uiServer;
        })
    }

    return new Promise((resolve, reject) => {
        const cwd = path.resolve(app.getAppPath(), '../ui');
        console.log(`>> Starting UI development server with command "npm run serve" in "${cwd}"...`);

        let resolved = false;
        const cmd = spawn('npm', ['run', 'serve'], {
            cwd,
            shell: true,
        });
        cmd.on('close', (code) => {
            console.log(`>> ZTF server closed with code ${code}`);
            _uiServerApp = null;
        });
        cmd.stdout.on('data', data => {
            if (resolved) {
                return;
            }
            const dataString = String(data);
            const lines = dataString.split('\n');
            for (let i = 0; i < lines.length; i++) {
                const line = lines[i];
                if (DEBUG) {
                    console.log('\t', line);
                }
                if (line.includes('App running at:')) {
                    const nextLine = lines[i + 1] || lines[i + 2];
                    if (DEBUG) {
                        console.log('\t', nextLine);
                    }
                    if (!nextLine) {
                        console.error('\t', `Cannot grabing running address after line "${line}".`);
                        throw new Error(`Cannot grabing running address after line "${line}".`);
                    }
                    const url = nextLine.split('Local:   ')[1];
                    if (url) {
                        resolved = true;
                        resolve(url);
                    }
                    if (!DEBUG) {
                        break;
                    }
                }
            }
        });
        cmd.on('error', spawnError => {
            console.error('>>> Get ui server url failed with error', spawnError);
            reject(spawnError)
        });
        _uiServerApp = cmd;
    });
}
