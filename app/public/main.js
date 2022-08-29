const btn = document.getElementById('btn')
const btnSolve = document.getElementById('btnSolve')
const solution = document.getElementById('solution')
const chars = document.getElementById('chars')
const numChars = document.getElementById('numChars')


// make an api call to https://bananagrams-sam.uc.r.appspot.com
// and generate a string of numChars characters
// and display it in the DOM
btn.addEventListener('click', () => {
    body = JSON.stringify({
        numChars: parseInt(numChars.value)
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
        chars.value = data
    })
    .catch(err => console.log(err))
    
})

btnSolve.addEventListener('click', () => {
    solution.innerText = "loading..."
    body = JSON.stringify({
        chars: chars.value
        })
    if (chars.value.length == 0) {
        solution.innerText = "no characters entered, try again"
        return
    }
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
        solution.innerText = data
    })
    .catch(err => console.log(err))
    
})