import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IstioComponent } from './istio.component';

describe('IstioComponent', () => {
  let component: IstioComponent;
  let fixture: ComponentFixture<IstioComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IstioComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IstioComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
