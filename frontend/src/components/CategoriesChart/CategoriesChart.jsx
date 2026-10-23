import { useEffect, useState } from 'react';
import styles from './CategoriesChart.module.css';

const CURRENCY_ICONS = {
  Рубли: '₽',
  Доллары: '$',
  Дирхамы: 'د.إ',
};

const formatNumber = (n) =>
  n.toLocaleString('ru-RU', { maximumFractionDigits: 2 });

function CategoriesChart() {
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch('/api/8')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => setRows(data ?? [])) // Go отдаёт null для пустого слайса
      .catch((err) => {
        console.error('Ошибка при загрузке', err);
        setError('Не удалось загрузить данные');
      })
      .finally(() => setLoading(false));
  }, []);

  // Для доли используем рублёвый эквивалент, если бэкенд его отдаёт
  const getRub = (r) => r.total_rub ?? r.total;
  const totalRub = rows.reduce((sum, r) => sum + getRub(r), 0);

  const categoriesData = rows
    .map((r, i) => ({
      id: `${r.currency}-${r.type_money}-${i}`,
      name: `${r.currency} · ${r.type_money}`,
      value: formatNumber(r.total),
      rub: getRub(r),
      percentage: totalRub > 0 ? (getRub(r) / totalRub) * 100 : 0,
      color: '#388e3c',
      icon: CURRENCY_ICONS[r.currency] ?? '💱',
    }))
    .sort((a, b) => b.rub - a.rub);

  if (loading) return <div className={styles.container}>Загрузка...</div>;
  if (error) return <div className={styles.container}>{error}</div>;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h3 className={styles.title}>Денежные средства</h3>
      </div>
      <div className={styles.chartArea}>
        {categoriesData.map((category) => (
          <div key={category.id} className={styles.categoryRow}>
            <div className={styles.categoryInfo}>
              <span className={styles.icon}>{category.icon}</span>
              <span className={styles.name}>{category.name}</span>
            </div>
            <div className={styles.valueAndBar}>
              <span className={styles.value}>{category.value}</span>
              <span className={styles.percentage}>
                {Math.round(category.percentage)}%
              </span>
              <div className={styles.progressBarBg}>
                <div
                  className={styles.progressBarFill}
                  style={{
                    width: `${category.percentage}%`,
                    backgroundColor: category.color,
                  }}
                />
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default CategoriesChart;