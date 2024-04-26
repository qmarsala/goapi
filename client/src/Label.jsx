import React, { useState, useEffect } from 'react';

function Label() {
    const [labels, setLabels] = useState([]);
    const [textInput, setTextInput] = useState('');
    const [targetInput, setTargetInput] = useState('');
    const [editingId, setEditingId] = useState(null);
    const [editText, setEditText] = useState('');
    const [editTarget, setEditTarget] = useState('');

    useEffect(() => {
        fetch('http://localhost:8080/api/labels')
            .then(response => response.json())
            .then(data => setLabels(data.labels))
            .catch(error => console.error('Error fetching labels:', error));
    }, []);

    const addLabel = () => {
        const newLabel = { text: textInput, target: targetInput };
        fetch('http://localhost:8080/api/labels', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newLabel)
        })
            .then(response => response.json())
            .then(label => {
                setLabels([...labels, label]);
                console.log(label)
                setTextInput('');
                setTargetInput('');
            })
            .catch(error => console.error('Error adding label:', error));
    };

    const startEditing = (label) => {
        console.log(label)
        setEditingId(label.id);
        setEditText(label.message);
    };

    const cancelEditing = () => {
        setEditingId(null);
        setEditText('');
    };

    const updateLabel = id => {
        const label = labels.find(l => l.id === id);
        fetch(`http://localhost:8080/api/labels/${label.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: label.id, message: editText })
        })
            .then(response => response.json())
            .then(updatedLabel => {
                const newlabels = labels.map(p => p.id === id ? updatedLabel : p);
                setLabels(newlabels);
                cancelEditing();
            })
            .catch(error => console.error('Error toggling label:', error));
    };

    const deleteLabel = id => {
        fetch(`http://localhost:8080/api/labels/${id}`, {
            method: 'DELETE'
        })
            .then(() => {
                setLabels(labels.filter(label => label.id !== id));
            })
            .catch(error => console.error('Error deleting label:', error));
    };

    return (
        <div>
            <input value={textInput} onChange={(e) => setTextInput(e.target.value)} />
            <input value={targetInput} onChange={(e) => setTargetInput(e.target.value)} />
            <button onClick={addLabel}>Add Label</button>
            <ul>
                {labels.map(label => (
                    <li key={label.id}>
                        {editingId === label.id ? (
                            <>
                                <input type="text" value={editText} onChange={(e) => setEditText(e.target.value)} />
                                <input type="text" value={editTarget} onChange={(e) => setEditTarget(e.target.value)} />
                                <button onClick={() => updateLabel(label.id)}>Save</button>
                                <button onClick={cancelEditing}>Cancel</button>
                            </>
                        ) : (
                            <>
                                {label.text}
                                {label.target}
                                <button onClick={() => startEditing(label)}>Edit</button>
                                <button onClick={() => deleteLabel(label.id)}>Delete</button>
                            </>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default Label;
