import './wails-bridge';
import './monaco';
import App from './views/App.svelte';
import { Events } from './runtime';

Events.Once('wombat:init', (data) => {
  const buildMode = data?.build_mode ?? data?.buildMode;
  if (buildMode === 'prod') {
    window.addEventListener('contextmenu', (e) => e.preventDefault());
  }
});

window.isWin = window.navigator.platform.startsWith('Win');

const app = new App({
  target: document.body,
});

export default app;
