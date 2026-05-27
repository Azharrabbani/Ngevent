export const toggleItem = <T>(
  items: T[],
  value: T
): T[] => {
  if (items.includes(value)) {
    return items.filter((item) => item !== value);
  }

  return [...items, value];
};