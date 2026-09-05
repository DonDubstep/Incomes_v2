import styles from './TransactionsList.module.css';

function TransactionsList() {
  // Данные транзакций с картинки-макета
  const transactions = [
    {
      id: 1,
      title: 'Зарплата',
      amount: '+120 000 ₽',
      date: '31 мая 2024',
      icon: '🕑',
    },
    {
      id: 2,
      title: 'Кешбек',
      amount: '+25 000 ₽',
      date: '28 мая 2024',
      icon: '%',
    },
    {
      id: 3,
      title: 'Инвестиции',
      subtitle: 'Дивиденды',
      amount: '+7 500 ₽',
      date: '27 мая 2024',
      icon: '📈',
    },
  ];

  return (
    <div className={styles.container}>
      <h3 className={styles.title}>Последние поступления</h3>
      
      <div className={styles.list}>
        {transactions.map((tx) => (
          <div key={tx.id} className={styles.row}>
            <div className={styles.leftGroup}>
              <div className={styles.iconCircle}>{tx.icon}</div>
              <div className={styles.info}>
                <div className={styles.txTitle}>{tx.title}</div>
              </div>
            </div>

            <div className={styles.rightGroup}>
              <div className={styles.amount}>{tx.amount}</div>
              <div className={styles.date}>{tx.date}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default TransactionsList;