import { define, html, render } from 'hybrids';


// XXX: Får ikke kjørt CORS preflight mot arbeidsplassen bruker proxy.
function hentStillinger()  {
    const url = "https://tovare.com/jobb/rss?view=json"
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


const SisteStillinger = {
        s: () => hentStillinger(),
        h: "hello",
        render: ({ s, h }) =>html `
        <style>


:host {
    display: block;
}
:host([hidden]) {
    display: none;
    box-sizing: border-box;
}

table {
    font-family: Arial, Helvetica, sans-serif;
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
                <th>Jobb</th><th>Arbeidsgiver</th><th>Sted</th>
            </thead>
            <tbody>
    ${content.map(({ title,employer,workLocations}) => html`
        <tr>
            <td>${title}</td>
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