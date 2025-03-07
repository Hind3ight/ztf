import {app, BrowserWindow} from 'electron';
import {getUIServerUrl, startZtfServer} from './services';

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) { // eslint-disable-line global-require
  app.quit();
}

const createWindow = (url) => {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 800,
    height: 600,
  });

  // and load the index.html of the app.
  mainWindow.loadURL(url);

  // Open the DevTools.
  mainWindow.webContents.openDevTools();
};


let _starting = false;

async function startApp() {
  if (_starting) {
    return;
  }
  _starting = true;

  try {
    const ztfServerUrl = await startZtfServer();
    console.log(`>> ZTF Server started successfully: ${ztfServerUrl}`);
  } catch (error) {
    console.error('>> Start ztf server failed: ' + error);
    process.exit(1);
    return;
  }

  const url = await getUIServerUrl();

  console.log('>> UI server url is', url);

  createWindow(url);

  _starting = false;
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', startApp);

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (BrowserWindow.getAllWindows().length === 0) {
    startApp();
  }
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and import them here.
