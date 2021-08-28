window.onload = () => {
  document.getElementById("subscribe").addEventListener("click", () => {
    if (!("serviceWorker" in navigator)) {
      alert("your browser does not support service worker");
      return;
    }
    if (Notification.permission !== "denied") {
      Notification.requestPermission().then(() => {
        navigator.serviceWorker.register("sw.js").then((swReg) => {
          fetch("/api/getPublicKey").then((data) => {
            data.text().then((key) => {
              console.log(key);
              console.log(swReg);
              swReg.pushManager
                .subscribe({
                  userVisibleOnly: true,
                  applicationServerKey: urlB64ToUint8Array(key),
                })
                .then((subscription) => {
                  console.log(subscription);
                  fetch("/api/subscribe", {
                    method: "POST",
                    body: JSON.stringify(subscription),
                  }).then(() => {
                    alert("success!");
                  });
                });
            });
          });
        });
      });
    }
  });

  function urlB64ToUint8Array(base64String) {
    const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
    const base64 = (base64String + padding)
      .replace(/\-/g, "+")
      .replace(/_/g, "/");

    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);

    for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
  }
};
