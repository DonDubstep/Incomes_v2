import { useEffect, useState } from 'react';
import styles from './TransactionsList.module.css';

const formatAmount = (n) =>
  '+' + n.toLocaleString('ru-RU', { maximumFractionDigits: 2 }) + ' ₽';

const formatDate = (str) => {
  // "2024-05-31" или "2024-05-31 12:00:00" -> "31 мая 2024"
  const d = new Date(str.replace(' ', 'T'));
  if (Number.isNaN(d.getTime())) return str; // если формат неожиданный, показываем как есть
  return d.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

function TransactionsList() {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch('/api/9')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => setTransactions(data ?? [])) // Go отдаёт null для пустого слайса
      .catch((err) => {
        console.error('Ошибка при загрузке', err);
        setError('Не удалось загрузить данные');
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className={styles.container}>Загрузка...</div>;
  if (error) return <div className={styles.container}>{error}</div>;

  return (
    <div className={styles.container}>
      <h3 className={styles.title}>Последние поступления</h3>
      <div className={styles.list}>
        {transactions.length === 0 && <div>Нет поступлений</div>}
        {transactions.map((tx, i) => (
          <div key={`${tx.date}-${i}`} className={styles.row}>
            <div className={styles.leftGroup}>
              <div className={styles.iconCircle}>💰</div>
              <div className={styles.info}>
                <div className={styles.txTitle}>{tx.description}</div>
              </div>
            </div>

            <div className={styles.rightGroup}>
              <div className={styles.amount}>{formatAmount(tx.amount)}</div>
              <div className={styles.date}>{formatDate(tx.date)}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default TransactionsList;