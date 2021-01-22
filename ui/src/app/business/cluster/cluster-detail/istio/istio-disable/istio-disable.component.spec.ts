import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IstioDisableComponent } from './istio-disable.component';

describe('IstioDisableComponent', () => {
  let component: IstioDisableComponent;
  let fixture: ComponentFixture<IstioDisableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IstioDisableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IstioDisableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
