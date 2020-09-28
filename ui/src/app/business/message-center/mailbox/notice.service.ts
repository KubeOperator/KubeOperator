import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {Observable} from 'rxjs';
import {Notice} from './notice';
import {Page} from '../../../shared/class/Page';

@Injectable({
    providedIn: 'root'
})
export class NoticeService extends BaseModelService<any> {

    baseUrl = '/api/v1/message/mail';
    constructor(http: HttpClient) {
        super(http);
    }

    pageBy(page, size, userName): Observable<Page<Notice>> {
        const pageUrl = `${this.baseUrl}?pageNum=${page}&pageSize=${size}&userName=${userName}`;
        return this.http.get<Page<Notice>>(pageUrl);
    }
}
