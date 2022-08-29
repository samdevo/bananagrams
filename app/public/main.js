const btn = document.getElementById('btn')
const btnSolve = document.getElementById('btnSolve')

// make an api call to https://bananagrams-sam.uc.r.appspot.com
// and generate a string of numChars characters
// and display it in the DOM
btn.addEventListener('click', () => {
    body = JSON.stringify({
        numChars: parseInt(document.getElementById('numChars').value)
        })
    fetch('https://bananagrams-func-xn6nonedta-uc.a.run.app/generate', {
        method: 'POST',
        headers: {
            'Host': 'bananagrams-func-xn6nonedta-uc.a.run.app',
            'Accept': '*/*'
        },
        body: body,
    })
    .then(res => res.text())
    .then(data => {
        console.log(data)
        document.getElementById('chars').value = data
    })
    .catch(err => console.log(err))
    
})

btnSolve.addEventListener('click', () => {
    document.getElementById('solution').innerText = "loading..."
    body = JSON.stringify({
        chars: document.getElementById('chars').value
        })
    fetch('https://bananagrams-func-xn6nonedta-uc.a.run.app/solve', {
        method: 'POST',
        headers: {
            'Host': 'bananagrams-func-xn6nonedta-uc.a.run.app',
            'Accept': '*/*'
        },
        body: body,
    })
    .then(res => res.text())
    .then(data => {
        console.log(data)
        document.getElementById('solution').innerText = data
    })
    .catch(err => console.log(err))
    
})