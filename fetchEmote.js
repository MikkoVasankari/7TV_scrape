import puppeteer from "puppeteer";

const arg = process.argv.slice(2);
let emotesFound = false;

(async () => {
  const browser = await puppeteer.launch({
    args: [
      "--no-sandbox",
      "--disable-setuid-sandbox",
      "--disable-dev-shm-usage",
    ],
    ignoreHTTPSErrors: true,
    headless: true,
  });

  const page = await browser.newPage();

  const emote = arg[0];
  await page.goto("https://7tv.app/emotes?page=1&query=" + emote, {
    waitUntil: "networkidle2",
  });

  try {
    await page.waitForSelector(".emote-card-wrapper", { visible: true, timeout: 10000 });
    emotesFound = true;
  } catch (error) {
    console.log("No emotes found");
  }

  const data = await page.evaluate(() => {
    const emoteCards = document.querySelectorAll(
      ".emote-card-wrapper .emote-card"
    );
    const result = [];

    emoteCards.forEach((card) => {
      const link = card.querySelector("a");
      const titleBanner = card.querySelector(".title-banner");

      if (link && titleBanner) {
        const href = link.getAttribute("href");
        const title = titleBanner.textContent.trim();
        result.push({ href, title });
      }
    });

    return result;
  });

  if (emotesFound == true) {
    console.log(JSON.stringify(data));
  }

  await browser.close();
})();
