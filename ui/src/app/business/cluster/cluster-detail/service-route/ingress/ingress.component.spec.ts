import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IngressComponent } from './ingress.component';

describe('IngressComponent', () => {
  let component: IngressComponent;
  let fixture: ComponentFixture<IngressComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IngressComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IngressComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
