const btn = document.getElementById('btn')
const btnSolve = document.getElementById('btnSolve')
const solution = document.getElementById('solution')
const chars = document.getElementById('chars')
const numChars = document.getElementById('numChars')
const loading = document.getElementById('loading')
const minLength = document.getElementById('minLength')
const maxLength = document.getElementById('maxLength')


// // automatically update the width of chars to match the input width
// chars.addEventListener('input', () => {
//     chars.style.width = chars.value.length * 10 + 'px'
// })


// make an api call to https://bananagrams-sam.uc.r.appspot.com
// and generate a string of numChars characters
// and display it in the DOM
btn.addEventListener('click', () => {
    // loading.innerText = "loading..."
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
        // loading.innerText = ""
    })
    .catch(err => console.log(err))
    
})

btnSolve.addEventListener('click', () => {
    solution.innerText = "loading..."
    body = JSON.stringify({
        chars: chars.value,
        minChars: parseInt(minLength.value),
        maxChars: parseInt(maxLength.value)
        })
    console.log(body)
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