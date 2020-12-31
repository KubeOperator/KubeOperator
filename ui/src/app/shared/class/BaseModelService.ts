import {BaseModel, BaseRequest} from './BaseModel';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Page} from './Page';
import {Batch} from './Batch';

export abstract class BaseModelService<T extends BaseModel> {

    baseUrl = '';
    variable = new Map<string, string>();

    protected constructor(protected http: HttpClient) {
    }

    list(): Observable<Page<T>> {
        const url = this.urlHandler();
        return this.http.get<Page<T>>(url);
    }

    page(page, size): Observable<Page<T>> {
        const url = this.urlHandler();
        const pageUrl = `${url}?pageNum=${page}&pageSize=${size}`;
        return this.http.get<Page<T>>(pageUrl);
    }

    get(name: string): Observable<T> {
        const url = this.urlHandler();
        const itemUrl = `${url}/${name}`;
        return this.http.get<T>(itemUrl);
    }

    create(item: BaseRequest): Observable<T> {
        const url = this.urlHandler();
        return this.http.post<T>(url, item);
    }

    update(name: string, item: BaseRequest): Observable<T> {
        const url = this.urlHandler();
        const itemUrl = `${url}/${name}/`;
        return this.http.patch<T>(itemUrl, item);
    }

    delete(name: string): Observable<any> {
        const url = this.urlHandler();
        const itemUrl = `${url}/${name}/`;
        return this.http.delete<any>(itemUrl);
    }

    batch(method: string, items: T[]): Observable<any> {
        const url = this.urlHandler();
        const batchUrl = `${url}/batch/`;
        const b = new Batch<T>(method, items);
        return this.http.post(batchUrl, b);
    }

    private urlHandler() {
        let url = this.baseUrl;
        this.variable.forEach(((k, v) => {
            if (url.indexOf(`{${k}}`) !== -1) {
                url = url.replace(`{${k}}`, v);
            }
        }));
        return url;
    }
}

