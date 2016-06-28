(function () {
    require.config({ paths: { 'vs': 'monaco-editor/min/vs' } });
    require(['vs/editor/editor.main'], function () {
        var editor = monaco.editor.create(document.getElementById('container'), {
            value: [
                'function x() {',
                '\tconsole.log("Hello world!");',
                '}'
            ].join('\n'),
            language: 'javascript'
        });

        //load in our d.ts file
        var xhr = new XMLHttpRequest();
        xhr.addEventListener("load", function () {
            monaco.languages.typescript.javascriptDefaults.addExtraLib(this.responseText, 'nano-service.d.ts');
        });
        xhr.open("GET", "nano-service.d.ts");
        xhr.send();

    });
} ());