/*
Copyright 2015, 2019 Google Inc. All Rights Reserved.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Incrementing OFFLINE_VERSION will kick off the install event and force
// previously cached resources to be updated from the network.
const PREFIX = "V1";
// Customize this with a different URL if needed.
const OFFLINE_URL = '/';

const FONT_AWESOME_CDN_FONT = 'https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.2.0/webfonts/fa-solid-900.woff2';
const FONT_AWESOME_CDN_CSS = 'https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.2.0/css/all.min.css';

self.addEventListener('install', (event) => {
  self.skipWaiting();

  event.waitUntil((async () => {
    const cache = await caches.open(PREFIX);

    await Promise.all([
      cache.add(new Request(OFFLINE_URL, {cache: 'reload'})),
      cache.add(new Request(FONT_AWESOME_CDN_FONT, {cache: 'reload'})),
      cache.add(new Request(FONT_AWESOME_CDN_CSS, {cache: 'reload'})),
    ]);
  })());

  console.log(`${PREFIX} install`)
});

self.addEventListener('activate', (event) => {
  console.log(`${PREFIX} active`)

  event.waitUntil((async () => {
    // Enable navigation preload if it's supported.
    // See https://developers.google.com/web/updates/2017/02/navigation-preload
    if ('navigationPreload' in self.registration) {
      await self.registration.navigationPreload.enable();
    }
  })());

  // Tell the active service worker to take control of the page immediately.
  clients.claim();

  event.waitUntil((async () => {
    const keys = await caches.keys();
    await Promise.all(keys.map(k => {
      if (k.includes(PREFIX)) {
        return caches.delete(k);
      }
    }))
  })())
});

self.addEventListener('fetch', (event) => {
  //console.log(`${PREFIX} Fetching : ${event.request.url}, Mode : ${event.request.mode}`);
  
  /*We only want to call event.respondWith() if this is a navigation request
  for an HTML page.*/
  if (event.request.mode === 'navigate') {
    event.respondWith((async () => {
      try {
        // First, try to use the navigation preload response if it's supported.
        /*const preloadResponse = await event.preloadResponse;
        if (preloadResponse) {
          return preloadResponse;
        }*/

        const cache = await caches.open(PREFIX);
        cache.add(event.request, {cache: 'reload'})

        const networkResponse = await fetch(event.request);
        return networkResponse;
      } catch (error) {
        /*catch is only triggered if an exception is thrown, which is likely
        due to a network error.
        If fetch() returns a valid HTTP response with a response code in
        the 4xx or 5xx range, the catch() will NOT be called.*/
        console.log('Fetch failed; returning offline page instead.', error);

        const cache = await caches.open(PREFIX);
        return await cache.match(OFFLINE_URL);
      }
    })());
  } else {
    event.respondWith((async () => {
      const cache = await caches.open(PREFIX);
      if (event.request.url.startsWith("{{.Hostname}}")) {
        if (event.request.url.indexOf('{{.Hostname}}/oauth') !== -1) {
          return await fetch(event.request);
        }

        try {
          // First, try to use the navigation preload response if it's supported.
          const preloadResponse = await event.preloadResponse;
          if (preloadResponse) {
            return preloadResponse;
          }
          
          cache.add(event.request, {cache: 'reload'})
  
          return await fetch(event.request);
        } catch (error) {
          /*catch is only triggered if an exception is thrown, which is likely
          due to a network error.
          If fetch() returns a valid HTTP response with a response code in
          the 4xx or 5xx range, the catch() will NOT be called.*/
          //console.log('Fetch failed; returning offline page instead.', error);
  
          const cache = await caches.open(PREFIX);
          return await cache.match(event.request.url);
        }
      } else if ((await cache.keys()).map(r => r.url).indexOf(event.request.url) !== -1) {
        return await cache.match(event.request.url);
      }

      //console.log(`Not navigate and not current domain : ${event.request.url}`);

      return await fetch(event.request);
    })())
  }

  /*If our if() condition is false, then this fetch handler won't intercept the
  request. If there are any other fetch handlers registered, they will get a
  chance to call event.respondWith(). If no fetch handlers call
  event.respondWith(), the request will be handled by the browser as if there
  were no service worker involvement.*/
});