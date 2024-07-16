(async () => {
  try {
    const [tab] = await chrome.tabs.query({
      active: true,
      lastFocusedWindow: true,
    });

    const videoIdMatch = tab.url.match(/(?<=v=)[^&]+/);
    if (videoIdMatch && videoIdMatch[0]) {
      const videoId = videoIdMatch[0];
      getTranscript(videoId);
    } else {
      throw new Error("Video ID not found");
    }
  } catch (error) {
    document.getElementById("summary").innerText = "";
    document.getElementById("error").innerText = error.message;
  }
})();

async function getTranscript(videoId) {
  try {
    const response = await fetch(`http://localhost:8080/api/transcript/${videoId}`);
    
    if (!response.ok) {
      throw new Error("Failed to get transcript: " + await response.text());
    }

    const data = await response.json();
    document.getElementById("summary").innerText = data.summary;
    document.getElementById("error").innerText = ""; // Clear any previous error message
  } catch (error) {
    document.getElementById("summary").innerText = "";
    document.getElementById("error").innerText = error.message;
  }
}