self.addEventListener("push", (event) => {
  let data = event.data.json();
  console.log("[sw] recv data", data);

  event.waitUntil(
    self.registration.showNotification(data.Title, {
      body: data.Body,
      data,
    })
  );
});

self.addEventListener("notificationclick", function (event) {
  if (event.notification.data && event.notification.data.URL !== "") {
    event.notification.close();
    console.log("[sw] notificationclick", event);
    clients.openWindow(event.notification.data.URL);
  }
});
