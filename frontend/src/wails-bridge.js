import * as API from '../wailsjs/go/app/API';

// Expose Go bindings for the Svelte UI. Do NOT overwrite window.wails —
// Wails v2 owns that object for EventsNotify / window flags.
window.backend = { api: API };
