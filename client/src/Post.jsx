import React, { useState, useEffect } from 'react';

function Post() {
    const [Posts, setPosts] = useState([]);
    const [input, setInput] = useState('');
    const [editingId, setEditingId] = useState(null);
    const [editText, setEditText] = useState('');

    useEffect(() => {
        fetch('http://localhost:8080/api/posts')
            .then(response => response.json())
            .then(data => setPosts(data.posts))
            .catch(error => console.error('Error fetching Posts:', error));
    }, []);

    const addPost = () => {
        const newPost = { message: input };
        fetch('http://localhost:8080/api/posts', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newPost)
        })
            .then(response => response.json())
            .then(post => {
                setPosts([...Posts, post]);
                setInput('');
            })
            .catch(error => console.error('Error adding post:', error));
    };

    const startEditing = (post) => {
        console.log(post)
        setEditingId(post.id);
        setEditText(post.message);
    };

    const cancelEditing = () => {
        setEditingId(null);
        setEditText('');
    };

    const updatePost = id => {
        const post = Posts.find(post => post.id === id);
        fetch(`http://localhost:8080/api/posts/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ...post })
        })
            .then(response => response.json())
            .then(updatedPost => {
                const newPosts = Posts.map(p => p.id === id ? updatedPost : p);
                setPosts(newPosts);
                cancelEditing();
            })
            .catch(error => console.error('Error toggling post:', error));
    };

    const deletePost = id => {
        fetch(`http://localhost:8080/api/posts/${id}`, {
            method: 'DELETE'
        })
            .then(() => {
                setPosts(Posts.filter(post => post.id !== id));
            })
            .catch(error => console.error('Error deleting post:', error));
    };

    return (
        <div>
            <input value={input} onChange={(e) => setInput(e.target.value)} />
            <button onClick={addPost}>Add Post</button>
            <ul>
                {Posts.map(post => (
                    <li key={post.id}>
                        {editingId === post.id ? (
                            <>
                                <input type="text" value={editText} onChange={(e) => setEditText(e.target.value)} />
                                <button onClick={() => updatePost(post.id)}>Save</button>
                                <button onClick={cancelEditing}>Cancel</button>
                            </>
                        ) : (
                            <>
                                {post.message}
                                <button onClick={() => startEditing(post)}>Edit</button>
                                <button onClick={() => deletePost(post.id)}>Delete</button>
                            </>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default Post;
