import { useEffect, useState } from 'react';
import styles from './StatsCards.module.css';

function StatsCards() {
  const [stats, setStats] = useState([]);
  const [loading, setLoading] = useState(true); 

  useEffect(() => {

    Promise.all([
      fetch('/api/5').then((res) => res.json()),
      fetch('/api/6').then((res) => res.json()),
      fetch('/api/2?period=2026-05').then((res) => res.json()),
    ])
    .then(([totalData, monthlyData, changeData]) => {
      setStats([
        {id: 1, title: 'Общая сумма', value: totalData, isMain:true },
        {id: 2, title: 'Доход за последний месяц', value: monthlyData, isMain:false },
        {id: 3, title: 'Изменение общей суммы', value: changeData, isMain:false },
      ]);
      setLoading(false);
    })
    .catch((error) => {
      console.error('Ошибка при загрузке', error);
      setLoading(false);
    });
  }, []);

  if(loading) {
    return <div>Загрузка...</div>;
  }

  // Находим главную карточку для левой колонки
  const mainCard = stats.find((item) => item.isMain);

  // Оставляем остальные карточки для правой колонки
  const secondaryCards = stats.filter((item) => !item.isMain);

  return (
    <div className={styles.container}>
      {/* 1. Главная карточка слева */}
      {mainCard && (
        <div className={`${styles.card} ${styles.mainCard}`}>
          <div className={styles.title}>{mainCard.title}</div>
          <div className={styles.largeValue}>{mainCard.value} ₽</div>
        </div>
      )}

      {/* 2. Правая колонка — создаём карточки через .map() */}
      <div className={styles.rightColumn}>
        {secondaryCards.map((card) => (
          <div key={card.id} className={styles.card}>
            <div className={styles.title}>{card.title}</div>
            <div className={styles.value}>{card.value} ₽</div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default StatsCards;