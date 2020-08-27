(function(document, window) {
    let q = document.getElementById('q'),
        options = document.getElementById('option'),
        submit = document.getElementById('search'),
        search_result = document.getElementById('search-results');
    /*
     * Fetch and load our document index
     */
    let requestURL = '/documents.json',
        request = new XMLHttpRequest();

    request.open('GET', requestURL);
    request.responseType = 'json';
    request.send();
    request.onload = function() {
        const data = request.response;
        setupSearchBox(data);
    };

    function populateResult(container, result, data) {
        console.log("data ->", data);
        let child = document.createElement('section'),
            src = [];
        src = [ 
            '<h3>', '<a href="/', result.ref, '/">', data.title, '</a>', '</h3>',
            '<div>', '<span class="label">', 'Interviewer:', '</span>', data.interviewer, '</div>',
            '<div>', '<span class="label">', 'Interview dates:', '</span>', data.interviewdate, '</div>',
            '<p>', '<span class="label">', 'Abstract:', '</span>', data['abstract'], '<p>'
        ].join('')
        console.log(src);
        child.innerHTML = src; 
        container.appendChild(child);
    }

    function buildResult(result) {
         console.log('buildResult: ' + result);
         let container = document.getElementById('result-list');
         let request = new XMLHttpRequest();
         request.open('GET', ['/', result.ref, '/scheme.json'].join(''));
         request.responseType = 'json';
         request.send();
         request.onload = function() {
             const data = request.response;
             console.log('data laoded for ref ' + result.ref);
             populateResult(container, result, data); 
         }
    }

    
    /*
     * This setups the search box to use Lunr.js search engine
     * and documents.json index.
     */
    function setupSearchBox(data) {
        let idx = lunr.Index.load(data);

        submit.addEventListener('click', function(event) {
            let terms = '';
            terms = q.value;
            console.log("Terms: "+terms);
            /* NOTE: we need to reset our results container! */
            search_result.innerHTML = '<section id="result-list"></section>';
            results = idx.search(terms);
            results.forEach(buildResult);
            event.preventDefault();
        }, false);
    }
}(document, window));
