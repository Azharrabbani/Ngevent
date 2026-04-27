import parsePhoneNumber from 'libphonenumber-js'


export function GetIsoFromPhoneNumber(phoneNumberString: string) {
    try {
        const phoneNumber = parsePhoneNumber(phoneNumberString)

        if (!phoneNumber || !phoneNumber.isValid()) {
            return undefined;
        }

        return phoneNumber.country;
    } catch {
        return undefined;
    }
}