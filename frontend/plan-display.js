(function () {
  const monetaryFields = new Set([
    "startingCapital",
    "monthlyContribution",
    "totalInterestEarnings",
    "totalEarnings",
    "totalDeposited",
    "interest",
    "tax",
    "contribution",
    "increase",
    "capital",
    "startingPrincipal",
    "monthlyPayment",
    "escrowPayment",
    "totalExpenditure",
    "totalPaid",
    "amount",
    "payment",
    "otherPayments",
    "paydown",
    "principal",
  ]);
  const dateFields = new Set(["date", "startDate"]);
  const rateFields = new Set([
    "yearlyInterestRate",
    "taxRate",
    "yearlyInflationRate",
    "monthlyInterestRate",
    "rateOfReturn",
    "inflationAdjustedROR",
    "costOfCredit",
    "costOfCreditPercent",
  ]);
  const fieldLabels = {
    inflationAdjustedROR: "Inflation Adjusted Rate of Return",
    escrowPayment: "Other Payments",
    totalEarnings: "Total Interest Earnings",
    costOfCreditPercent: "Cost Of Credit",
  };

  function humanize(value) {
    if (fieldLabels[value]) {
      return fieldLabels[value];
    }

    return value
      .replace(/([a-z0-9])([A-Z])/g, "$1 $2")
      .replace(/[_-]/g, " ")
      .replace(/^./, (letter) => letter.toUpperCase());
  }

  function formatFieldValue(field, value) {
    if (value === null || value === undefined) {
      return "Not provided";
    }

    if (monetaryFields.has(field)) {
      return formatters.dollars(value);
    }

    if (dateFields.has(field)) {
      return formatters.date(value);
    }

    if (rateFields.has(field)) {
      return formatters.percentage(value);
    }

    return String(value);
  }

  function createDetailLine(labelText, valueText) {
    const line = document.createElement("div");
    line.className =
      "grid gap-1 border-b border-base-300 px-5 py-4 last:border-b-0 sm:grid-cols-[minmax(12rem,1fr)_2fr] sm:gap-6";

    const label = document.createElement("dt");
    label.className = "font-medium text-base-content/65";
    label.textContent = labelText;

    const value = document.createElement("dd");
    value.className = "break-words";
    value.textContent = valueText === null ? "null" : String(valueText);

    line.append(label, value);
    return line;
  }

  function createDefinitionList(data) {
    const list = document.createElement("dl");
    list.className =
      "overflow-hidden rounded-box border border-base-300 bg-base-100 shadow-sm";

    Object.entries(data).forEach(([key, value]) => {
      list.appendChild(createDetailLine(humanize(key), formatFieldValue(key, value)));
    });

    return list;
  }

  function createDataTable(rows, emptyMessage) {
    if (!Array.isArray(rows) || rows.length === 0) {
      const emptyState = document.createElement("div");
      emptyState.className = "alert";
      emptyState.textContent = emptyMessage;
      return emptyState;
    }

    const columns = [...new Set(rows.flatMap((row) => Object.keys(row)))];
    const wrapper = document.createElement("div");
    wrapper.className =
      "overflow-x-auto rounded-box border border-base-300 bg-base-100 shadow-sm";
    const table = document.createElement("table");
    table.className = "table table-zebra";

    const head = document.createElement("thead");
    const headerRow = document.createElement("tr");
    columns.forEach((column) => {
      const header = document.createElement("th");
      header.textContent = humanize(column);
      headerRow.appendChild(header);
    });
    head.appendChild(headerRow);

    const body = document.createElement("tbody");
    rows.forEach((row) => {
      const tableRow = document.createElement("tr");
      columns.forEach((column) => {
        const cell = document.createElement("td");
        cell.textContent = formatFieldValue(column, row[column]);
        tableRow.appendChild(cell);
      });
      body.appendChild(tableRow);
    });

    table.append(head, body);
    wrapper.appendChild(table);
    return wrapper;
  }

  function createSection(titleText, content) {
    const section = document.createElement("section");
    const title = document.createElement("h2");
    title.className = "mb-4 text-2xl font-bold";
    title.textContent = titleText;
    section.append(title, content);
    return section;
  }

  function renderResults(container, sections) {
    const fragment = document.createDocumentFragment();
    sections.forEach((section) => {
      fragment.appendChild(createSection(section.title, section.content));
    });
    container.replaceChildren(fragment);
    container.classList.remove("hidden");
  }

  window.planDisplay = Object.freeze({
    renderSavingsResults(container, data) {
      const { plan, ...calculated } = data || {};
      renderResults(container, [
        {
          title: "Calculated Information",
          content: createDefinitionList(calculated),
        },
        {
          title: "Savings Plan",
          content: createDataTable(
            plan,
            "No savings plan data was returned by the backend.",
          ),
        },
      ]);
    },
    renderLoanResults(container, data) {
      const { plan, costOfCreditPercent, durationMonths, ...rest } = data || {};
      const calculated = {
        durationMonths: Number.isFinite(Number(durationMonths))
          ? `${durationMonths} months`
          : durationMonths,
        ...rest,
        costOfCredit: costOfCreditPercent,
      };
      renderResults(container, [
        {
          title: "Calculated Information",
          content: createDefinitionList(calculated),
        },
        {
          title: "Payment Plan",
          content: createDataTable(
            plan,
            "No payment plan data was returned by the backend.",
          ),
        },
      ]);
    },
  });
})();
