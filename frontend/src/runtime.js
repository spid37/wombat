import {
  EventsOn,
  EventsEmit,
  EventsOnce,
  BrowserOpenURL,
} from '../wailsjs/runtime/runtime';

// App-facing runtime helpers (replaces Wails v1's global `wails` API).
export const Events = {
  On: EventsOn,
  Emit: EventsEmit,
  Once: EventsOnce,
};

export const Browser = {
  OpenURL: BrowserOpenURL,
};
