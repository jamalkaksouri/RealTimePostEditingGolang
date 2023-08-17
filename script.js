document.addEventListener("DOMContentLoaded", () => {
  const eventSource = new EventSource("http://localhost:6835/sse");

  let currentStockValue = "";
  const container = document.querySelector("#stock-element");

  eventSource.addEventListener("message", (event) => {
    const eventData = JSON.parse(event.data);

    if (eventData.data === ":keepalive") {
      return;
    }

    if (eventData.event === "time") {
      const timeElement = document.getElementById("time-element");
      timeElement.textContent = "Current Time: " + eventData.data;
    } else if (eventData.event === "isStock") {
      const stockElement = document.getElementById("stock-element");

      if (eventData.data.startsWith("In stock:")) {
        const newValue = extractNumericValue(eventData.data);

        if (currentStockValue !== newValue) {
          currentStockValue = newValue;
          animateContainer2(container);
        }
      } else if (eventData.data === "Product X is out of stock") {
        animateContainer(container);
      }
      stockElement.textContent = eventData.data;
    }
  });

  function extractNumericValue(str) {
    const matches = str.match(/\d+/);
    return matches ? matches[0] : "";
  }

  function animateContainer(container) {
    container.style.animation = "swashOut 2s ease-in-out";
  }

  function animateContainer2(container) {
    container.style.animation = "puffIn 1s ease-in-out";
    setTimeout(() => {
      container.style.animation = ""; // Remove animation after it completes
    }, 1000);
  }

  // Rest of the script remains the same
});
