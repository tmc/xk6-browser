import { chromium } from 'k6/x/browser';

export const options = {
  scenarios: {
    ui: {
      executor: 'shared-iterations',
      options: {
        browser: {
            type: 'chromium',
        },
      },
    },
  },
  thresholds: {
    checks: ["rate==1.0"]
  }
}

/*
Page Object Model is a well-known pattern to abstract a web page.

The Locator API enables using the Page Object Model pattern to organize
and simplify test code.

Note: For comparison, you can see another example that does not use
the Page Object Model pattern in locator.js.
*/
export class Bet {
  constructor(page) {
    this.page = page;
    this.headsButton = page.locator("input[value='Bet on heads!']");
    this.tailsButton = page.locator("input[value='Bet on tails!']");
    this.currentBet = page.locator("//p[starts-with(text(),'Your bet: ')]");
  }

  goto() {
    return this.page.goto("https://test.k6.io/flip_coin.php", { waitUntil: "networkidle" });
  }

  heads() {
    return Promise.all([
      this.page.waitForNavigation(),
      this.headsButton.click(),
    ]);
  }

  tails() {
    return Promise.all([
      this.page.waitForNavigation(),
      this.tailsButton.click(),
    ]);
  }

  current() {
    return this.currentBet.innerText();
  }
}

export default async function() {
  const browser = chromium.launch();
  const context = browser.newContext();
  const page = context.newPage();

  const bet = new Bet(page);
  try {
    await bet.goto()
    await bet.tails();
    console.log("Current bet:", bet.current());
    await bet.heads();
    console.log("Current bet:", bet.current());
    await bet.tails();
    console.log("Current bet:", bet.current());
    await bet.heads();
    console.log("Current bet:", bet.current());
  } finally {
    page.close();
    browser.close();
  }
}
