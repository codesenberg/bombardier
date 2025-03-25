document.addEventListener('DOMContentLoaded', function() {
    // DOM elements
    const bombardierForm = document.getElementById('bombardierForm');
    const testTypeRequests = document.getElementById('testTypeRequests');
    const testTypeDuration = document.getElementById('testTypeDuration');
    const requestsInput = document.getElementById('requestsInput');
    const durationInput = document.getElementById('durationInput');
    const testStatus = document.getElementById('testStatus');
    const testResults = document.getElementById('testResults');
    const testError = document.getElementById('testError');
    const saveConfigBtn = document.getElementById('saveConfig');
    const loadConfigBtn = document.getElementById('loadConfig');
    const configModal = new bootstrap.Modal(document.getElementById('configModal'));
    const saveConfigModal = new bootstrap.Modal(document.getElementById('saveConfigModal'));
    const confirmSaveConfig = document.getElementById('confirmSaveConfig');

    // Toggle between requests/duration
    testTypeRequests.addEventListener('change', function() {
        if (this.checked) {
            requestsInput.classList.remove('d-none');
            durationInput.classList.add('d-none');
        }
    });

    testTypeDuration.addEventListener('change', function() {
        if (this.checked) {
            durationInput.classList.remove('d-none');
            requestsInput.classList.add('d-none');
        }
    });

    // Form submission
    bombardierForm.addEventListener('submit', function(e) {
        e.preventDefault();
        runTest();
    });

    // Save configuration
    saveConfigBtn.addEventListener('click', function() {
        saveConfigModal.show();
    });

    // Confirm save configuration
    confirmSaveConfig.addEventListener('click', function() {
        const configName = document.getElementById('configName').value.trim();
        if (configName === '') {
            alert('Please enter a name for this configuration.');
            return;
        }

        const config = getFormData();
        
        // Save to localStorage
        let savedConfigs = JSON.parse(localStorage.getItem('bombardierConfigs') || '{}');
        savedConfigs[configName] = config;
        localStorage.setItem('bombardierConfigs', JSON.stringify(savedConfigs));
        
        saveConfigModal.hide();
        alert('Configuration saved successfully!');
    });

    // Load configuration
    loadConfigBtn.addEventListener('click', function() {
        const savedConfigs = JSON.parse(localStorage.getItem('bombardierConfigs') || '{}');
        const configNames = Object.keys(savedConfigs);
        
        if (configNames.length === 0) {
            alert('No saved configurations found.');
            return;
        }
        
        // Populate the modal
        const savedConfigsList = document.getElementById('savedConfigs');
        savedConfigsList.innerHTML = '';
        
        configNames.forEach(name => {
            const item = document.createElement('button');
            item.className = 'list-group-item list-group-item-action';
            item.textContent = name;
            
            item.addEventListener('click', function() {
                loadConfig(savedConfigs[name]);
                configModal.hide();
            });
            
            savedConfigsList.appendChild(item);
        });
        
        configModal.show();
    });

    // Run test
    function runTest() {
        // Show status and clear previous results/errors
        testStatus.classList.remove('d-none');
        testResults.textContent = '';
        testError.classList.add('d-none');
        testError.textContent = '';

        // Get form data
        const data = getFormData();
        
        // Send API request
        fetch('/api/run', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
            testStatus.classList.add('d-none');
            
            if (data.success) {
                testResults.textContent = data.output;
            } else {
                testError.classList.remove('d-none');
                testError.textContent = data.errorMessage;
            }
        })
        .catch(error => {
            testStatus.classList.add('d-none');
            testError.classList.remove('d-none');
            testError.textContent = 'Error: ' + error.message;
        });
    }

    // Get form data
    function getFormData() {
        // Parse headers from textarea
        const headersText = document.getElementById('headers').value;
        const headers = {};
        
        if (headersText.trim() !== '') {
            headersText.split('\n').forEach(line => {
                if (line.trim() !== '') {
                    const firstColon = line.indexOf(':');
                    if (firstColon > 0) {
                        const key = line.substring(0, firstColon).trim();
                        const value = line.substring(firstColon + 1).trim();
                        headers[key] = value;
                    }
                }
            });
        }

        // Create the request object
        const data = {
            url: document.getElementById('url').value,
            method: document.getElementById('method').value,
            numConnections: parseInt(document.getElementById('numConnections').value),
            headers: headers,
            body: document.getElementById('body').value,
            timeout: parseInt(document.getElementById('timeout').value),
            latencies: document.getElementById('latencies').checked,
            insecureSkipVerify: document.getElementById('insecureSkipVerify').checked,
            disableKeepAlives: document.getElementById('disableKeepAlives').checked,
            clientType: document.getElementById('clientType').value
        };

        // Add rate limit if specified
        const rateInput = document.getElementById('rate').value;
        if (rateInput && rateInput.trim() !== '') {
            data.rate = parseInt(rateInput);
        }

        // Add test type (requests or duration)
        if (document.getElementById('testTypeRequests').checked) {
            data.numRequests = parseInt(document.getElementById('numRequests').value);
        } else {
            data.duration = parseInt(document.getElementById('duration').value);
        }

        return data;
    }

    // Load config into form
    function loadConfig(config) {
        document.getElementById('url').value = config.url;
        document.getElementById('method').value = config.method;
        document.getElementById('numConnections').value = config.numConnections;
        
        // Test type
        if (config.numRequests !== undefined) {
            document.getElementById('testTypeRequests').checked = true;
            document.getElementById('numRequests').value = config.numRequests;
            requestsInput.classList.remove('d-none');
            durationInput.classList.add('d-none');
        } else if (config.duration !== undefined) {
            document.getElementById('testTypeDuration').checked = true;
            document.getElementById('duration').value = config.duration;
            durationInput.classList.remove('d-none');
            requestsInput.classList.add('d-none');
        }
        
        document.getElementById('timeout').value = config.timeout;
        
        // Rate
        if (config.rate !== undefined) {
            document.getElementById('rate').value = config.rate;
        } else {
            document.getElementById('rate').value = '';
        }
        
        document.getElementById('latencies').checked = config.latencies;
        document.getElementById('insecureSkipVerify').checked = config.insecureSkipVerify;
        document.getElementById('disableKeepAlives').checked = config.disableKeepAlives;
        document.getElementById('clientType').value = config.clientType;
        
        // Headers
        let headersText = '';
        for (const key in config.headers) {
            headersText += `${key}: ${config.headers[key]}\n`;
        }
        document.getElementById('headers').value = headersText.trim();
        
        document.getElementById('body').value = config.body || '';
    }
}); 