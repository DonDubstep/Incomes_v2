// src/components/StatsCards/StatsCards.jsx
import styles from './StatsCards.module.css';

function StatsCards() {
  const statsData = [
    {
      id: 1,
      title: 'Общий доход',
      value: '245 120 ₽',
      trend: '+12.5%',
      trendText: 'по сравнению с апрелем',
      large: true
    },
    {
      id: 2,
      title: 'Средний доход в день',
      value: '7 907 ₽',
    },
    {
      id: 3,
      title: 'Количество транзакций',
      value: '28',
    }
  ];

  return (
    <div className={styles.statsGrid}>
      {statsData.map((stat) => (
        <div 
          key={stat.id} 
          className={`${styles.card} ${stat.large ? styles.largeCard : ''}`}
        >
          <div className={styles.title}>{stat.title}</div>
          <div className={styles.value}>{stat.value}</div>
          
          {stat.trend && (
            <div className={styles.trend}>
              <span className={styles.trendArrow}>↑</span>
              <span className={styles.trendValue}>{stat.trend}</span>
              <span className={styles.trendText}>{stat.trendText}</span>
            </div>
          )}
        </div>
      ))}
    </div>
  );
}

export default StatsCards;