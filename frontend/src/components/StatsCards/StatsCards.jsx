import styles from './StatsCards.module.css';

// 📦 Тестовые данные прямо в JSX-файле
const initialStats = [
  {
    id: 1,
    title: 'Общая сумма',
    value: '280 000',
    isMain: true,
  },
  {
    id: 2,
    title: 'Доход за последний месяц',
    value: '7 907',
    isMain: false,
  },
  {
    id: 3,
    title: 'Изменение общей суммы',
    value: '28',
    isMain: false,
  },
];

function StatsCards() {
  // Находим главную карточку для левой колонки 👈
  const mainCard = initialStats.find((item) => item.isMain);

  // Оставляем остальные карточки для правой колонки 👉
  const secondaryCards = initialStats.filter((item) => !item.isMain);

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