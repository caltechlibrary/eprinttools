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
            fields = [],
            src = '';
        fields = [ 
            '<h3>', '<a href="/', result.ref, '/">', data.title, '</a>', '</h3>'
        ]
        if (data.type) {
            fields.push('<div>', '<span class="label">', 'Document Type:', '</span>', data.type, '</div>');
        }
        if (data.date) {
            fields.push('<div>', '<span class="label">', 'Date:', '</span>', data.date, '</div>');
        } else if (data.year) {
            fields.push('<div>', '<span class="label">', 'Year:', '</span>', data.year, '</div>');
        }
        if (data.doi) {
            fields.push('<div>', '<span class="label">', 'DOI:', '</span>', data.doi, '</div>');
        }
        if (data.interviewer) {
            fields.push('<div>', '<span class="label">', 'Interviewer:', '</span>', data.interviewer, '</div>');
        }
        if (data.interviewdate) {
            fields.push('<div>', '<span class="label">', 'Interview dates:', '</span>', data.interviewdate, '</div>');
        }
        if (data['abstract']) {
            fields.push('<p>', '<span class="label">', 'Abstract:', '</span>', data['abstract'], '<p>')
        }
        if (data.collection) {
            fields.push('<p>', '<span class="label">', 'Collection:', '</span>', data.collection, '<p>')
        }
        src = fields.join('')
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
