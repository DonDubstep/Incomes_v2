import { useState } from 'react';

function BalanceCard({title, initialBalance}) {
// Используем переданный initialBalance как начальное значение состояния
  const [balance, setBalance] = useState(initialBalance);

  return (
    <div style={{ 
      border: '1px solid #ccc', 
      borderRadius: '8px', 
      padding: '16px', 
      marginBottom: '16px' 
    }}>
      <h3>{title} 💳</h3>
      <h2>Баланс: {balance} ₽</h2>
      <button onClick={() => setBalance(balance + 100)}>
        Пополнить на 100 ₽
      </button>
      {' '}
      <button onClick={() => setBalance(balance - 100)}>
        Списать 100 ₽
      </button>
    </div>
  );
}

export default BalanceCard;