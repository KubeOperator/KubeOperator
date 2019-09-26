import {StorageCredential} from './storage-credential';

export class BackupStorage {
    id: string;
    name: string;
    region: string;
    credentials: StorageCredential;
    type: string;
    status: string;
    date_created: string;
}
