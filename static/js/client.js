'use strict';

async function wake(macAddr, name = "") {
    const displayName = name != "" ? name + " (" + macAddr + ")" : macAddr

    const button = document.getElementById(macAddr + ".Button");
    if (button) {
        button.disabled = true;
        button.innerText = "Waking...";
    }

    try {
        const response = await fetch("/api/v1/wake/" + macAddr);

        const responseBody = await response.json();

        if (response.ok) {
            if (button) {
                button.innerText = "✅ Woken up";
                // Wait 1 second before reverting the button text
                await new Promise(resolve => setTimeout(resolve, 1000));
            } else {
                appendAlert("Send magic packet to " + displayName);
            }
        } else {
            if (button) {
                button.innerText = "❌ Failed";
                // Wait 1 second before reverting the button text
                await new Promise(resolve => setTimeout(resolve, 1000));
            } else {
                appendAlert("Failed to send magic packet to " + displayName + " : " + responseBody.reason, "warning");
            }
        }
    } catch (error) {
        console.error(error.message);
        if (button) {
            button.innerText = "❌ Failed";
            // Wait 1 second before reverting the button text
            await new Promise(resolve => setTimeout(resolve, 1000));
        } else {
            appendAlert("Failed to send magic packet to " + displayName, "danger");
        }
    } finally {
        if (button) {
            button.innerText = "Wake";
            button.disabled = false;
        }
    }
}

async function wakeCustom() {
    const inputCustomMAC = document.getElementById("custom-mac-input");
    wake(inputCustomMAC.value);
}

async function addHost() {
    const name = document.getElementById('hostName').value;
    const macAddr = document.getElementById('macAddress').value;
    const address = document.getElementById('address').value;

    const host = {
        mac: macAddr,
        name: name
    }
    if (address != "") {
        host.address = address;
    }

    modal.hide();
    try {
        const response = await fetch('/api/v1/hosts', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(host)
        });

        const responseBody = await response.json();

        if (response.ok) {
            appendAlert(`Added host ${name}`);
            location.reload();
        } else {
            appendAlert(`Failed to add host: ${responseBody.reason}`, "warning");
        }
    } catch (error) {
        console.error(error.message);
        appendAlert("Failed to add host " + name, "danger");
    }
}

async function deleteHost(macAddr, name) {
    if (!confirm(`Are you sure you want to delete ${name}?`)) {
        return;
    }

    try {
        const response = await fetch(`/api/v1/hosts/${macAddr}`, {
            method: 'DELETE',
        });

        const responseBody = await response.json();

        if (response.ok) {
            appendAlert(`Deleted host ${name}`);
            location.reload();
        } else {
            appendAlert(`Failed to delete host: ${responseBody.reason}`, "warning");
        }
    } catch (error) {
        console.error(error.message);
        appendAlert("Failed to delete host " + name, "danger");
    }
}

async function hostStatus() {
    try {
        const response = await fetch(`/api/v1/hosts/status`);

        const responseBody = await response.json();

        if (response.ok) {
            updateHostStatus(responseBody);
        } else {
            appendAlert(`Failed to fetch host status: ${responseBody.reason}`, "warning");
        }
    } catch (error) {
        console.error(error.message);
        appendAlert("Failed to fetch host status", "danger");
    }
}
