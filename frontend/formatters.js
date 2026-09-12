(function () {
  const dollarFormatter = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

  function dollars(amountInCents) {
    const amount = Number(amountInCents);
    if (!Number.isFinite(amount)) {
      return "Not available";
    }

    return dollarFormatter.format(amount / 100);
  }

  function date(value) {
    const text = String(value);
    const datePrefix = text.match(/^\d{4}-\d{2}-\d{2}/);
    if (datePrefix) {
      return datePrefix[0];
    }

    const parsedDate = new Date(value);
    if (Number.isNaN(parsedDate.getTime())) {
      return text;
    }

    return parsedDate.toISOString().slice(0, 10);
  }

  function percentage(value) {
    const rate = Number(value);
    if (!Number.isFinite(rate)) {
      return "Not available";
    }

    return `${rate.toLocaleString("en-US", {
      maximumFractionDigits: 10,
    })}%`;
  }

  window.formatters = Object.freeze({
    dollars,
    date,
    percentage,
  });
})();
