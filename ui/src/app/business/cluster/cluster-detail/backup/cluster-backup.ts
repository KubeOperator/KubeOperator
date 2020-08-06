import {BaseModel} from '../../../../shared/class/BaseModel';

export class BackupStrategy extends BaseModel {
    id: string;
    cron: number;
    saveNum: number;
    backupAccountName: string;
    status: string;
    clusterName: string;
}

export class BackupFile extends BaseModel {
    name: string;
    clusterName: string;
    clusterBackupStrategyID: string;
    folder: string;
}