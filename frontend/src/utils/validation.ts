export const mapValidationErrors = (error: any) => {
    const validationError = error?.response?.data?.error;

    if (!Array.isArray(validationError)) return {};

    return validationError.reduce((acc: any, e: any) => {
        acc[e.field] = e.message;
        return acc;
    }, {});
};