import { define, html, render,property } from 'hybrids';


// Får ikke kjørt CORS preflight mot arbeidsplassen bruker proxy.
function hentStillinger(antall)  {
    const url = "https://tovare.com/jobb/rss?view=json&size="+antall
    const req = {
        method: 'GET',
        mode: 'cors',
        cache: 'default'
    };
    return fetch(url, req)
        .then(data => {
            // We promise some json
            return data.json()
        })
};

// Component
const SisteStillinger = {
        antall: 5,
        s: ({ antall }) => hentStillinger(antall),
        h: "hello",
        render: ({ s, h }) =>html `
<style>
:host {
    display: block;
    box-sizing: border-box;
}
:host([hidden]) {
    display: none;
    box-sizing: border-box;
}

table {
    font-family: Arial, Helvetica, sans-serif;
    font-size: 14px;
}

th, td {
  padding: 15px;
  text-align: left;
  border-bottom: 1px solid #ddd;
} 
</style>
    ${html.resolve(
        s.then(({totalElements, content}) => html`
        <table>
            <thead>
                <th>Jobb</th><th>Arbeidsgiver</th><th>Kommune</th>
            </thead>
            <tbody>
    ${content.map(({ title,link,employer,workLocations}) => html`
        <tr>
            <td><a href="${link}">${title}</a></td>
            <td>${employer.name}</td>
            <td>${workLocations[0].municipal}</td>
        </tr>`)}        
            </tbody>
        `)
          .catch(() => "ERROR"), html`
          <p><h1 style="text-align: center;">Henter stillinger...</h1></p>`,1000)
      }
      </table>
      <p style="text-align: center;">Kilde: Ledige stillinger fra <a href="https://arbeidsplassen.no">arbeidsplassen</a></p>
    `
}
define('siste-stillinger', SisteStillinger)
