import {Injectable} from '@angular/core';

@Injectable({
    providedIn: 'root'
})
export class BusinessLicenseService {
    licenseValid = false;

    constructor() {
    }

    update(licenseValid: boolean) {
        this.licenseValid = licenseValid;
    }

    get selectedValue() {
        return this.licenseValid;
    }
}
