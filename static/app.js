const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const req = `/search?q=${data.query}&work=${data.work}`;
    
    const response = fetch(req)
    .then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    table.innerHTML = rows.join(" ");
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
