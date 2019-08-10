import React from 'react';

const Mes = ({chat, user}) => (
    <li className={`chat ${user == chat.username ? "right" : "left"}`}>
        {chat.content}
    </li>
);

export default Mes;