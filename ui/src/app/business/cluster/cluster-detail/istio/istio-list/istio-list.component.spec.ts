import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IstioListComponent } from './istio-list.component';

describe('IstioListComponent', () => {
  let component: IstioListComponent;
  let fixture: ComponentFixture<IstioListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IstioListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IstioListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
