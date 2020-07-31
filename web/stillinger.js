import { define, html, render,property } from 'hybrids';

function hentStillinger(token,antall,municipal,orgnr,county)  {
    const auth = "Bearer "+token
    const url = "https://arbeidsplassen.nav.no/public-feed/api/v1/ads?size="+antall
        + "&municipal=" + municipal
        + "&orgnr=" + orgnr
        + "&county=" + county
    const req = {
        method: 'GET',
        withCredentials: true,
        credentials: 'include',
        headers: new Headers({
            'Authorization': auth,
            'Host': 'arbeidsplassen.nav.no'
          })
    };
    return fetch(url, req)
        .then(data => {
            // We promise some json
            return data.json()
        })
};

// Component
const Stillinger = {
        antall: 5,
        municipal: "",
        orgnr: "",
        county: "",
        token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJwdWJsaWMudG9rZW4udjFAbmF2Lm5vIiwiYXVkIjoiZmVlZC1hcGktdjEiLCJpc3MiOiJuYXYubm8iLCJpYXQiOjE1NTc0NzM0MjJ9.jNGlLUF9HxoHo5JrQNMkweLj_91bgk97ZebLdfx3_UQ",
        s: ({ token, antall, municipal,orgnr, county}) => hentStillinger(token, antall,municipal,orgnr,county),
        h: "hello",
        render: ({ antall, s, h }) =>html `
<style>
</style>
    ${html.resolve(
        s.then(({totalElements, content}) => html`
    <div>
    ${content.map(({ title,link,employer,workLocations,description}) => html`
        <h1><a href="${link}">${title}</a></h1>
        <h2>${employer.name} / ${workLocations[0].municipal}</h2>
        <div innerHTML="${description}"></div>

        `)}        
    `).catch(() => "ERROR"), html`
          <p><h1 style="text-align: center;">Henter stillinger...</h1></p>`,1000)
      }
     </div>
     <p style="text-align: center;">Kilde: Ledige stillinger fra <a href="https://arbeidsplassen.no">arbeidsplassen</a></p>
    `
}
define('nav-stillinger', Stillinger)
