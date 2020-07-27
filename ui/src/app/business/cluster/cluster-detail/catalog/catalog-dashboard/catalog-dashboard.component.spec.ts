import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CatalogDashboardComponent } from './catalog-dashboard.component';

describe('CatelogDashboardComponent', () => {
  let component: CatalogDashboardComponent;
  let fixture: ComponentFixture<CatalogDashboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CatalogDashboardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CatalogDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
